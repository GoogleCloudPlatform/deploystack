cat go.mod| grep github.com/GoogleCloudPlatform/deploystack | head -1 | sed  's/github.com\/GoogleCloudPlatform\///g' | xargs  > dsVersion
cat go.mod| grep github.com/GoogleCloudPlatform/deploystack/tui  |  sed  's/github.com\/GoogleCloudPlatform\///g' | xargs > dsTUIVersion
cat go.mod| grep github.com/GoogleCloudPlatform/deploystack/gcloud  |  sed  's/github.com\/GoogleCloudPlatform\///g' | xargs > dsGcloudVersion
date +"%m-%d-%Y %r %Z" | xargs > buildTime