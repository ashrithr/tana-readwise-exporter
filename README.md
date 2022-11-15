# Readwise Exporter for Tana.io

`tana-readwise-exporter` uses Readwise Export API to fetch the highlights and formats the exported data using [Tana Paste](https://help.tana.inc/build-tutorials/tana-paste.html) format, which can be coppied to your Tana worksapce.

## Prerequisites

Generate a Readwise Access Token from [here](https://readwise.io/access_token). You'd need to pass this token as a command line option to the script.

## Installation

```bash
go get github.com/ashrithr/tana-readwise-exporter
```

## Examples

### Export highlights for last 2 days

```bash
tana-readwise-exporter export --token <your-readwise-token> --updated-after 2
```

> On MAC you can copy the output fo the above command directly to the clipboard by piping the above command to `pbcopy`.

### Fetching highlights for specific books

step 1: list the books and get the Readwise ID(s) for specific books

```bash
tana-readwise-exporter list --token <your-readwise-token> --category books
```

step 2: get the highlights for specific books id's obtained from previous command

```bash
tana-readwise-exporter export --token <your-readwise-token> --ids <replace-id1-from-prev-cmd>
```

> you can fetch highlights for multiple books by passing comma separated list of ids `--ids 1234567,2345678`
