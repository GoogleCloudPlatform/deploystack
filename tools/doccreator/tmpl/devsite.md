Project: /shell/docs/cloud-shell-tutorials/deploystack/_project.yaml
Book: /shell/docs/cloud-shell-tutorials/deploystack/_book.yaml
{% include "/shell/docs/cloud-shell-tutorials/deploystack/_local_variables.html" %}
{% set stack_key = "{{$.GetShortNameUnderScore}}" %}
{% set stack_name %}{{"{{"}} repos[stack_key].label {{"}}"}}{% endset %}
{% set stack_url %}{{"{{"}} repos[stack_key].url {{"}}"}}{% endset %}
{% set stack_diagram %}{{"{{"}} repos[stack_key].diagram {{"}}"}}{% endset %}
{% set stack_products = repos[stack_key].products  %}

<!-- TODO: Review, place iun the right spot and remove from file -->
<!-- 
 "{{$.GetShortNameUnderScore}}": {
      "url": "{{$.GitRepo | ToLower }}", 
      "label": "{{.DeployStack.Title | ToLower |  Title }}",
      "diagram" : "arch-{{$.GetShortName}}.svg",
      "products" : [{{range $val := .Products}}"{{$val.GetShortNameUnderScore}}",{{end}}]
  },

 -->

<<_template_diagram.md>>

<!-- TODO: FILL IN -->
# {{"{{"}} stack_name {{"}}"}}
{{"{{"}} stack_name {{"}}"}} [FILL in]

<<_template_getting_started.md>>

<<_template_products.md>>

<<_template_scripts.md>>

### `./main.tf`


{{range $val := .Entities}}{{if .IsResource}}
<!-- TODO: FILL IN -->
#### [FILL IN]
[FILL IN]

``` hcl
{{$val.Text | TrimSpace}}
```
{{end}}{{end}}

<hr class="full-width">

## Conclusion
<!-- TODO: FILL IN -->
Once run you should now have [FILL IN ]. Additionally you should have all of the code to modify or extend this solution to fit your environment.
