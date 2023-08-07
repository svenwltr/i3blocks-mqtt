# i3blocks-mqtt

Command for templating message for MQTT to use in [i3blocks](https://github.com/vivien/i3blocks).

> **Development Status:** Currently, I am only using this myself.

> This project is open source, but not open for contribution, since helping contributors to ship a PR usually takes more time than doing it myself. Please drop an issue, if you want something changed. Also of course it is possible to fork the whole project.

## Example

```
[temperature-inside]
label=In:
interval=persist
format=json
command=i3block-mqtt subscribe
BLOCKS_BROKER=mqtt://localhost:1883
BLOCKS_TOPIC=wadra/dashboard/temperature/wohnzimmer
BLOCKS_JSON=true
BLOCKS_TEMPLATE_FULL={{. | printf "%.1f"}}°C
BLOCKS_TEMPLATE_COLOR={{if gt . 24.}}#cc0000{{end}}{{if lt . .20}}#73d216{{end}}
```

Prints lines like this to stdout:

```
{"full_text":"27.6°C","color":"#cc0000"}
```

Which show `In:27.6°C` in a red color in the i3 bar.
