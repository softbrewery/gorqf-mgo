all:

clean:

test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace

test-watch:
	ginkgo watch -r -v

test-travis:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover -coverprofile=coverage.txt -covermode=atomic --trace --race --compilers=2

toc:
	# toc is generated by https://github.com/ekalinin/github-markdown-toc
	gh-md-toc --insert ./README.md
	rm -Rf README.md.*