# Geneanet Parser
geneparse downloads and parses Geneanet backups as if it was the Geneanet Android app.

## How do I build it?
```sh
$ make build
```

## Usage
```sh
$ ./geneparse                                                                                                                                    ✔  system  
A tool to download, extract and parse Geneanet bases.

Usage:
  geneparse [flags]
  geneparse [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  dlextr      download and extract Geneanet bases
  gedcom      parse Geneanet bases and create a gedcom file
  help        Help about any command

Flags:
  -h, --help   help for geneparse

Use "geneparse [command] --help" for more information about a command.

$ ./geneparse dlextr --help
The dlextr command will connect to Geneanet as if it was the Geneanet Android app, and will download the Geneanet bases. These bases use the Geneweb format.

Usage:
  geneparse dlextr [flags]

Flags:
  -h, --help               help for dlextr
  -o, --outputdir string   Output directory for Geneanet bases (default "output")
  -p, --password string    Password to log in to Geneanet (required)
  -t, --timeout string     Connection timeout for requests to Geneanet (default "10s")
  -u, --username string    Username or email address to log in to Geneanet (required)

$ ./geneparse gedcom --help                                                                                                                                                     ✔  system  
The gedcom command will parse Geneanet bases downloaded by the dlextr command and will create the corresponding gedcom file.

Usage:
  geneparse gedcom [flags]

Flags:
  -h, --help              help for gedcom
  -i, --inputdir string   Input directory for Geneanet bases (default "output")
```

## Usage example
```sh
$ ./geneparse dlextr -u user -p password -o outputdir
2021/12/17 14:43:17 Session cookie (gntsess) value: 0_user_deadbeedeadbeefdeadbeefdeadbeefd
2021/12/17 14:43:18 Account infos:
{
  "privilege": "deadbeefdeadbeefdeadbeefdeadbeef",
  "tree": 1,
  "otherTrees": [
    [
      "othertree",
      "The Other Tree"
    ],
    [
      "johndoe",
      "John Doe"
    ]
  ],
  "tree_access": 1,
  "annivAsc": 4,
  "annivDesc": 3,
  "proprio": 100,
  "login": "user",
  "canEdit": 1,
  "quotaPictureExceeded": 0
}
2021/12/17 14:43:21 Processing file: pb_base_family.dat
2021/12/17 14:43:21 Processing file: pb_base_info.dat
2021/12/17 14:43:21 Processing file: pb_base_person_note.dat
2021/12/17 14:43:21 Processing file: pb_base_person_note.inx
2021/12/17 14:43:21 Processing file: pb_base_family_note.dat
2021/12/17 14:43:21 Processing file: pb_base_family_note.inx
2021/12/17 14:43:21 Processing file: pb_base_ascends.inx
2021/12/17 14:43:21 Processing file: pb_base_name.inx
2021/12/17 14:43:21 Processing file: pb_base_name.wi
2021/12/17 14:43:21 Processing file: pb_base_person.inx
2021/12/17 14:43:21 Processing file: pb_base_family.inx
2021/12/17 14:43:21 Processing file: pb_base_name.i
2021/12/17 14:43:21 Processing file: pb_base_name.w
2021/12/17 14:43:21 Processing file: pb_base_person.dat
```
