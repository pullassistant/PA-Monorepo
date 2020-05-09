#!/bin/sh

REACT_APP_API_URL=https://app.pullassistant.com/api
REACT_APP_GA_TRACKER=UA-145781573-1
npm run build
rm -r ./../service/kodata/dashboard/*
mv -v build/* ../service/kodata/dashboard/