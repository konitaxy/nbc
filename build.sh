GOOS=linux GOARCH=amd64 go build -o build/server
cp config_pro.yaml build
cp config_dev.yaml build
cp -r resource build
cp -r build /Users/bohu/project/p-card/deploy

cd /Users/bohu/project/p-card/deploy
git add .
git commit -m 'backend update'
git push --set-upstream origin main