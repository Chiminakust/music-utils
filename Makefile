##
# Playlist downloader
#
# @file
# @version 0.1

CHAPTER_SPLITTER_PKG = ./cmd/chapter-splitter
BUILDDIR = ./build

all: chapter_splitter

chapter_splitter:
	go build -o $(BUILDDIR)/$@ $(CHAPTER_SPLITTER_PKG)

clean:
	$(RM) $(BUILDDIR)/*

.PHONY: clean
# end
