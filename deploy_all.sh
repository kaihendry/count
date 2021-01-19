for b in .git/refs/heads/*
do
	bash ./trigger_build.sh ${b##*/}
done
