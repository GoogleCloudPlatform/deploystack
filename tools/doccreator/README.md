# Documentation Creator

Will stub out a Deploystack documentation page based on targeted github repo.

usage

```shell
cd tools/doccreator
./clean && go run *.go -repo [a public github repo with a deploystack project]

```

It will output an index.html file to ./out/[repo name]/index.html. Use that for 
inclusion at deploystack.dev.