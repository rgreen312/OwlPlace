#!/usr/bin/env bash

function update-pixel {
    local x="$1"
    local y="$2"
    local r="$3"
    local g="$4"
    local b="$5"

    curl "localhost:3001/update_pixel?X=$x&Y=$y&R=$r&G=$g&B=$b"
}

update-pixel 10 10 100 100 100
update-pixel 11 11 100 100 100
update-pixel 12 12 100 100 100
update-pixel 13 13 100 100 100
