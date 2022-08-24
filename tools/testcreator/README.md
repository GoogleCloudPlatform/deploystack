# Test Creator

Will stub out a Deploystack test based on targeted terraform directory.

usage

```shell
cd tools/testcreator
./clean && go run *.go -folder [root folder of a deplpoystack project]

```

It will output a test file to ./out. Move that to your deploystack project folder. 