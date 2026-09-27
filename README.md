# stem-share

Send WAV stems to someone over the internet from the command line.

Uses the [Magic Wormhole](https://magic-wormhole.readthedocs.io/) protocol: transfers are end-to-end encrypted, work through NAT, and need no server of your own.

## Build

    go build -o stemshare .

## Usage

    stemshare send kick.wav snare.wav bass.wav
    stemshare receive 7-guitarist-revenge [-o ./stems] [-force]

Only `.wav` files with a valid RIFF/WAVE header are sent and saved.
