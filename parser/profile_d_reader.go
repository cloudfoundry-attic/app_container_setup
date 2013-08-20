package parser

func generateProfileDReader() string {
	return `
unset GEM_PATH
if [ -d app/.profile.d ]; then
for i in app/.profile.d/*.sh; do
  if [ -r $i ]; then
	. $i
  fi
done
unset i
fi
`
}
