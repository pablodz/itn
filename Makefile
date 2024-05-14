
list:
	git pull --tags
	git for-each-ref --sort=creatordate --format='%(creatordate:short)	%(refname:short)' refs/tags | tail -n 3