FROM tarantool/tarantool

# copy init lua script inside
COPY init_tarantool.lua /opt/tarantool/

CMD ["tarantool", "/opt/tarantool/init_tarantool.lua"]
