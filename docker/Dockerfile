from busybox

WORKDIR /app

COPY ./server /app/server
COPY ./build /app/build


ENTRYPOINT ["/app/server" "-s" "./build"]