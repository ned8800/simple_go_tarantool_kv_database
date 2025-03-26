box.cfg {
  listen = '0.0.0.0:3301'
}

box.schema.user.create('storage_user', { password = 'passw0rd', if_not_exists = true })
box.schema.user.grant('storage_user', 'super', nil, nil, { if_not_exists = true })

require('msgpack').cfg { encode_invalid_as_nil = true }


local db = box.schema.space.create('json_kv_database', {
  format = {
    { name = 'id',   type = 'string' },
    { name = 'data', type = 'string' }
  },
  if_not_exists = true,
})

db:create_index('primary', {
  parts = { 'id' },
  if_not_exists = true
})


db:insert({ "1", '{"name": "John Doe", "age": 30}' })
