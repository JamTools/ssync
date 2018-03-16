# ssync

Synchronize files between two paths.

On subsequent use, it removes files intentionally deleted from one path on the other.

## Usage

```
Usage: ssync [OPTIONS] LABEL PATH1 PATH2

Positional Args:
  LABEL           give label for subsequent use
  PATH1           1st directory path
  PATH2           2nd directory path
```

## Process

1. Looks for hidden file '.ssync-LABEL' within PATH1 or PATH2 which represents last shared state (list of common paths).

2. Compares latest state with current paths within each root path (PATH1, PATH2) to determine:

   * what has been deleted and should be deleted on the opposite path, then deletes.

   * what is new and should be copied to the opposite path, then copies.

3. Compares modified timestamp of each common path to determine more recently modified, then updates on opposite path.

4. Saves shared state file with updated common paths.

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

Source updated profile, run:

    source ~/.profile

### Building

Get latest source, run:

    go get github.com/JamTools/ssync

Navigate to source path, run:

    cd $GOPATH/src/github.com/JamTools/ssync

From within source path, run:

    go build

The binary will build to the current directory. To test by displaying usage, run:

    ./ssync --help

### Submitting a Pull Request

Fork repo on Github.

From within source path, setup new remote, run:

    git remote add myfork git@github.com:$GITHUB-USERNAME/ssync.git

Create a new branch to use for development, run:

    git checkout -b new-branch

Make your changes, add, commit and push to your Github fork.

Back on Github, submit pull request.

## License

This code is available open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
