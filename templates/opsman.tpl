 {% include "headers.tpl" %}

<div class="container">
    <h1>PCF installation details</h1>
    <h5>Details about Installed Product *</h5>
    <table class="table table-hover fixed">
        <thead class="thead-inverse">
            <tr>
                <th>Installed Products</th>
                <th>Version</th>
                <th>Release Notes</th>
            </tr>
        </thead>
        <tbody>

            {% for product in infoPcf.TileResources %}
            {% if product.Name !="p-bosh" %}
            <tr>
                <td class="accordion-toggle" data-toggle="collapse" data-target="#collapse{{forloop.Counter}}">
                    <b>{{product.Name}}</b><a href=""> (details)</a>
                    <div id="collapse{{forloop.Counter}}" class="collapse in">
                        <br>
                        <ul>
                            <li>Descriptions: {{product.Release.Description}}</li>
                            <li>Release Date: {{product.Release.ReleaseDate}}</li>
                            <li>Release Type: {{product.Release.ReleaseType}}</li>
                            <li>End Of Support: {{product.Release.EndOfSupportDate}}</li>
                        </ul>
                    </div>
                </td>
                <td>{{product.CleanVersion}}</td>
                <td><a target="_blank" href="{{product.Release.ReleaseNotesURL}}#{{product.CleanVersion}}">Click</a></td>
            </tr>
            {% endif %}
            {% endfor %}
        </tbody>
    </table>
    * Information Updated every new deployment


<hr/>
<br/>
    <h5>Buildpacks Details *</h5>
    <table class="table table-hover fixed">
        <thead class="thead-inverse">
            <tr>
                <th>Installed Buildpacs</th>
                <th>Version</th>
                <th>FileName</th>
                <th>Release Notes</th>
            </tr>
        </thead>
        <tbody>

            {% for bp in infoBuildpacks sorted  %}
            <tr>
                <td><b>{{bp.Name}}</b></td>
                <td>{{bp.Version}}</td>
                <td>{{bp.Filename}}
                <td><a target="_blank" href="{{bp.ReleaseNotesUrl}}">Click Here</a></td>
            </tr>
            {% endfor %}
        </tbody>
    </table>Updated every minutes
</div>
{% include "footers.tpl" %}
