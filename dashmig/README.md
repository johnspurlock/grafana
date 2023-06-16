# dashmig

dashboard.json is a modified copy of dashboardv38.json - modified to fit the current schema format.

goal: add a dashboard schema to our existing lineage that fits the v38 dashboard schema, and write the lenses to move between that and the current schema.

Currenlty blocked on Validate() errors returned from what I _think_ are valid fields:

```
validation against latest schema failed: no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.annotations: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.editable: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.fiscalYearStartMonth: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.graphTooltip: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.links: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.liveNow: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.panels: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.refresh: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.schemaVersion: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.style: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.tags: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.templating: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.time: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.timepicker: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.timezone: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.title: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.uid: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.version: field not allowed
no Thema handler for CUE error, please file an issue against github.com/grafana/thema
to improve this error output!

lineage._sortedSchemas.0._#schema.weekStart: field not allowed
```
