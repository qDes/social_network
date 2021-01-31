 wrk -t1 -c1 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase' && \
  wrk -t1 -c10 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase' && \
   wrk -t1 -c100 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase' && \
    wrk -t1 -c1000 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase'