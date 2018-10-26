# ssync

Synchronize files between two paths by comparing modified timestamps.

On subsequent use, using a shared state file, remove files intentionally deleted from one path on the other.

This is useful for a large music collection shared among a group of people, where individuals periodically change the collection by finding a better source and replacing, editing id3 tags to provide more accurate data, embedding artwork within audio files, deleting low quality recordings, etc.

Individuals can sync with multiple other individuals. Every two individuals have a shared state file. The latest shared state propagates across all individuals over time.

## Usage

```
Usage: ssync [OPTIONS] LABEL PATH1 PATH2

Positional Args:
  LABEL           give label for subsequent use
  PATH1           1st directory path
  PATH2           2nd directory path

Options:
  -confirm
      confirm before deleting files
  -force int
      update modified using this path regardless of modified timestamp (0=timestamp, 1=PATH1, 2=PATH2)
  -version
      print program version, then exit
```

## Process

1. Looks for state file '.ssync-LABEL' within root of PATH1 and PATH2 which contains list of common paths.

2. Compares latest state with current paths within each of PATH1 and PATH2 to determine:

   * what has been deleted and should be deleted on the opposite path, then deletes.

   * what is new and should be copied to the opposite path, then copies.

3. Compares modified timestamp of each common path to determine more recently modified, then updates on opposite path.

4. Saves updated state to file on both paths.

## Developing

### Install Go on Linux

Download latest go binary from [golang.org/dl](https://golang.org/dl/). In this case, version 1.10.

Extract to `/usr/local`, run:

    sudo tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz

Open ~/.profile for editing, run:

    nano ~/.profile

Append the following, then save/exit:

    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin

Source updated profile and ensure go $HOME/go exists, run:

    source ~/.profile
    if [ ! -d $HOME/go ]; then mkdir $HOME/go; fi

### Building

Get latest source, run:

    go get github.com/jamlib/ssync

Navigate to source path, run:

    cd $GOPATH/src/github.com/jamlib/ssync

From within source path, run:

    go build

The binary will build to the current directory. To test by displaying usage, run:

    ./ssync --help

### Testing

From within source path, run:

    go test -cover -v ./...

### Submitting a Pull Request

Fork repo on Github.

From within source path, setup new remote, run:

    gituser='YOUR-GITHUB-USERNAME'
    git remote add myfork git@github.com:$gituser/ssync.git

Create a new branch to use for development, run:

    git checkout -b new-branch

Make your changes, add, commit and push to your Github fork.

Back on Github, submit pull request.

## License

This code is available open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
