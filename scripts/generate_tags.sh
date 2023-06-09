# Lib modules
modules=(
  "firebase"
  "hash"
  "log"
  "mail"
  "migrate/sql"
  "random"
  "responses"
  "shutdown"
  "subscription"
  "tokens"
)

read -rp 'Enter tag: ' tag

git tag "$tag"

for module in "${modules[@]}"
do
   git tag "$module/$tag"
done

git push --tags