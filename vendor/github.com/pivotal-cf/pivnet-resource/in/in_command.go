package in

import (
	"fmt"
	"path/filepath"
	"strings"

	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logger"
	"github.com/pivotal-cf/pivnet-resource/concourse"
	"github.com/pivotal-cf/pivnet-resource/metadata"
	"github.com/pivotal-cf/pivnet-resource/versions"
)

//go:generate counterfeiter --fake-name FakeFilter . filterer
type filterer interface {
	ProductFileKeysByGlobs(
		productFiles []pivnet.ProductFile,
		globs []string,
	) ([]pivnet.ProductFile, error)
}

//go:generate counterfeiter --fake-name FakeDownloader . downloader
type downloader interface {
	Download(productFiles []pivnet.ProductFile, productSlug string, releaseID int) ([]string, error)
}

//go:generate counterfeiter --fake-name FakeFileSummer . fileSummer
type fileSummer interface {
	SumFile(filepath string) (string, error)
}

//go:generate counterfeiter --fake-name FakeFileWriter . fileWriter
type fileWriter interface {
	WriteMetadataJSONFile(mdata metadata.Metadata) error
	WriteMetadataYAMLFile(mdata metadata.Metadata) error
	WriteVersionFile(versionWithFingerprint string) error
}

//go:generate counterfeiter --fake-name FakePivnetClient . pivnetClient
type pivnetClient interface {
	GetRelease(productSlug string, version string) (pivnet.Release, error)
	AcceptEULA(productSlug string, releaseID int) error
	FileGroupsForRelease(productSlug string, releaseID int) ([]pivnet.FileGroup, error)
	ProductFilesForRelease(productSlug string, releaseID int) ([]pivnet.ProductFile, error)
	ProductFileForRelease(productSlug string, releaseID int, productFileID int) (pivnet.ProductFile, error)
	ReleaseDependencies(productSlug string, releaseID int) ([]pivnet.ReleaseDependency, error)
	DependencySpecifiers(productSlug string, releaseID int) ([]pivnet.DependencySpecifier, error)
	ReleaseUpgradePaths(productSlug string, releaseID int) ([]pivnet.ReleaseUpgradePath, error)
}

type InCommand struct {
	logger       logger.Logger
	downloadDir  string
	pivnetClient pivnetClient
	filter       filterer
	downloader   downloader
	fileSummer   fileSummer
	fileWriter   fileWriter
}

func NewInCommand(
	logger logger.Logger,
	pivnetClient pivnetClient,
	filter filterer,
	downloader downloader,
	fileSummer fileSummer,
	fileWriter fileWriter,
) *InCommand {
	return &InCommand{
		logger:       logger,
		pivnetClient: pivnetClient,
		filter:       filter,
		downloader:   downloader,
		fileSummer:   fileSummer,
		fileWriter:   fileWriter,
	}
}

func (c *InCommand) Run(input concourse.InRequest) (concourse.InResponse, error) {
	productSlug := input.Source.ProductSlug

	version, fingerprint, err := versions.SplitIntoVersionAndFingerprint(input.Version.ProductVersion)
	if err != nil {
		c.logger.Info("Parsing of fingerprint failed; continuing without it")
		version = input.Version.ProductVersion
		fingerprint = ""
	}

	c.logger.Info(fmt.Sprintf(
		"Getting release for product slug: '%s' and product version: '%s'",
		productSlug,
		version,
	))

	release, err := c.pivnetClient.GetRelease(productSlug, version)
	if err != nil {
		return concourse.InResponse{}, err
	}

	if fingerprint != "" {
		actualFingerprint := release.UpdatedAt
		if actualFingerprint != fingerprint {
			return concourse.InResponse{}, fmt.Errorf(
				"provided fingerprint: '%s' does not match actual fingerprint (from pivnet): '%s' - %s",
				fingerprint,
				actualFingerprint,
				"pivnet does not support downloading old versions of a release",
			)
		}
	}

	c.logger.Info(fmt.Sprintf("Accepting EULA for release with ID: %d", release.ID))

	err = c.pivnetClient.AcceptEULA(productSlug, release.ID)
	if err != nil {
		return concourse.InResponse{}, err
	}

	c.logger.Info("Getting product files")

	releaseProductFiles, err := c.pivnetClient.ProductFilesForRelease(productSlug, release.ID)
	if err != nil {
		return concourse.InResponse{}, err
	}

	c.logger.Info("Getting file groups")

	fileGroups, err := c.pivnetClient.FileGroupsForRelease(productSlug, release.ID)
	if err != nil {
		return concourse.InResponse{}, err
	}

	allProductFiles := releaseProductFiles
	for _, fg := range fileGroups {
		allProductFiles = append(allProductFiles, fg.ProductFiles...)
	}

	// Get individual product files to obtain metadata that isn't found
	// in the endpoint for all product files.
	for i, p := range allProductFiles {
		allProductFiles[i], err = c.pivnetClient.ProductFileForRelease(
			productSlug,
			release.ID,
			p.ID,
		)
		if err != nil {
			return concourse.InResponse{}, err
		}
	}

	c.logger.Info("Getting release dependencies")

	releaseDependencies, err := c.pivnetClient.ReleaseDependencies(productSlug, release.ID)
	if err != nil {
		return concourse.InResponse{}, err
	}

	c.logger.Info("Getting dependency specifiers")

	dependencySpecifiers, err := c.pivnetClient.DependencySpecifiers(productSlug, release.ID)
	if err != nil {
		return concourse.InResponse{}, err
	}

	c.logger.Info("Getting release upgrade paths")

	releaseUpgradePaths, err := c.pivnetClient.ReleaseUpgradePaths(productSlug, release.ID)
	if err != nil {
		return concourse.InResponse{}, err
	}

	c.logger.Info("Downloading files")

	err = c.downloadFiles(input.Params.Globs, allProductFiles, productSlug, release.ID)
	if err != nil {
		return concourse.InResponse{}, err
	}

	c.logger.Info("Creating metadata")

	versionWithFingerprint, err := versions.CombineVersionAndFingerprint(version, fingerprint)

	mdata := metadata.Metadata{
		Release: &metadata.Release{
			ID:                    release.ID,
			Version:               release.Version,
			ReleaseType:           string(release.ReleaseType),
			ReleaseDate:           release.ReleaseDate,
			Description:           release.Description,
			ReleaseNotesURL:       release.ReleaseNotesURL,
			Availability:          release.Availability,
			Controlled:            release.Controlled,
			ECCN:                  release.ECCN,
			LicenseException:      release.LicenseException,
			EndOfSupportDate:      release.EndOfSupportDate,
			EndOfGuidanceDate:     release.EndOfGuidanceDate,
			EndOfAvailabilityDate: release.EndOfAvailabilityDate,
		},
	}

	if release.EULA != nil {
		mdata.Release.EULASlug = release.EULA.Slug
	}

	for _, pf := range releaseProductFiles {
		mdata.Release.ProductFiles = append(mdata.Release.ProductFiles, metadata.ReleaseProductFile{
			ID: pf.ID,
		})
	}

	for _, pf := range allProductFiles {
		mdata.ProductFiles = append(mdata.ProductFiles, metadata.ProductFile{
			ID:                 pf.ID,
			File:               pf.Name,
			Description:        pf.Description,
			AWSObjectKey:       pf.AWSObjectKey,
			FileType:           pf.FileType,
			FileVersion:        pf.FileVersion,
			MD5:                pf.MD5,
			DocsURL:            pf.DocsURL,
			SystemRequirements: pf.SystemRequirements,
		})
	}

	for _, d := range releaseDependencies {
		mdata.Dependencies = append(mdata.Dependencies, metadata.Dependency{
			Release: metadata.DependentRelease{
				ID:      d.Release.ID,
				Version: d.Release.Version,
				Product: metadata.Product{
					ID:   d.Release.Product.ID,
					Name: d.Release.Product.Name,
				},
			},
		})
	}

	for _, d := range dependencySpecifiers {
		mdata.DependencySpecifiers = append(mdata.DependencySpecifiers, metadata.DependencySpecifier{
			ID:          d.ID,
			Specifier:   d.Specifier,
			ProductSlug: d.Product.Slug,
		})
	}

	for _, d := range releaseUpgradePaths {
		mdata.UpgradePaths = append(mdata.UpgradePaths, metadata.UpgradePath{
			ID:      d.Release.ID,
			Version: d.Release.Version,
		})
	}

	for _, fg := range fileGroups {
		mfg := metadata.FileGroup{
			ID:   fg.ID,
			Name: fg.Name,
		}

		for _, pf := range fg.ProductFiles {
			mfg.ProductFiles = append(mfg.ProductFiles, metadata.FileGroupProductFile{
				ID: pf.ID,
			})
		}

		mdata.FileGroups = append(mdata.FileGroups, mfg)
	}

	c.logger.Info("Writing metadata files")

	err = c.fileWriter.WriteVersionFile(versionWithFingerprint)
	if err != nil {
		return concourse.InResponse{}, err
	}

	err = c.fileWriter.WriteMetadataYAMLFile(mdata)
	if err != nil {
		return concourse.InResponse{}, err
	}

	err = c.fileWriter.WriteMetadataJSONFile(mdata)
	if err != nil {
		return concourse.InResponse{}, err
	}

	concourseMetadata := c.addReleaseMetadata([]concourse.Metadata{}, release)

	out := concourse.InResponse{
		Version: concourse.Version{
			ProductVersion: versionWithFingerprint,
		},
		Metadata: concourseMetadata,
	}

	return out, nil
}

func (c InCommand) downloadFiles(
	globs []string,
	productFiles []pivnet.ProductFile,
	productSlug string,
	releaseID int,
) error {
	c.logger.Info("Filtering download links by glob")

	filtered := productFiles

	// If globs were not provided, download everything without filtering.
	if globs != nil {
		var err error
		filtered, err = c.filter.ProductFileKeysByGlobs(productFiles, globs)
		if err != nil {
			return err
		}
	}

	c.logger.Info("Downloading filtered files")

	files, err := c.downloader.Download(filtered, productSlug, releaseID)
	if err != nil {
		return err
	}

	fileMD5s := map[string]string{}
	for _, p := range productFiles {
		parts := strings.Split(p.AWSObjectKey, "/")

		if len(parts) < 1 {
			panic("not enough components to form filename")
		}

		fileName := parts[len(parts)-1]

		if fileName == "" {
			panic("empty file name")
		}

		if p.FileType == pivnet.FileTypeSoftware {
			fileMD5s[fileName] = p.MD5
		}
	}

	err = c.compareMD5s(files, fileMD5s)
	if err != nil {
		return err
	}

	return nil
}

func (c InCommand) addReleaseMetadata(
	concourseMetadata []concourse.Metadata,
	release pivnet.Release,
) []concourse.Metadata {
	cmdata := append(concourseMetadata,
		concourse.Metadata{Name: "version", Value: release.Version},
		concourse.Metadata{Name: "release_type", Value: string(release.ReleaseType)},
		concourse.Metadata{Name: "release_date", Value: release.ReleaseDate},
		concourse.Metadata{Name: "description", Value: release.Description},
		concourse.Metadata{Name: "release_notes_url", Value: release.ReleaseNotesURL},
		concourse.Metadata{Name: "availability", Value: release.Availability},
		concourse.Metadata{Name: "controlled", Value: fmt.Sprintf("%t", release.Controlled)},
		concourse.Metadata{Name: "eccn", Value: release.ECCN},
		concourse.Metadata{Name: "license_exception", Value: release.LicenseException},
		concourse.Metadata{Name: "end_of_support_date", Value: release.EndOfSupportDate},
		concourse.Metadata{Name: "end_of_guidance_date", Value: release.EndOfGuidanceDate},
		concourse.Metadata{Name: "end_of_availability_date", Value: release.EndOfAvailabilityDate},
	)

	if release.EULA != nil {
		concourseMetadata = append(concourseMetadata,
			concourse.Metadata{Name: "eula_slug", Value: release.EULA.Slug},
		)
	}

	return cmdata
}

func (c InCommand) compareMD5s(filepaths []string, expectedMD5s map[string]string) error {
	c.logger.Info("Calcuating MD5 for downloaded files")

	for _, downloadPath := range filepaths {
		_, f := filepath.Split(downloadPath)

		actualMD5, err := c.fileSummer.SumFile(downloadPath)
		if err != nil {
			return err
		}

		expectedMD5 := expectedMD5s[f]
		if expectedMD5 != "" && expectedMD5 != actualMD5 {
			return fmt.Errorf(
				"MD5 comparison failed for downloaded file: '%s'. Expected (from pivnet): '%s' - actual (from file): '%s'",
				downloadPath,
				expectedMD5,
				actualMD5,
			)
		}
	}

	c.logger.Info("MD5 matched for all downloaded files")

	c.logger.Info("Get complete")

	return nil
}
