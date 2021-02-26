# Proposed protocol for turn-like server **lori**

- Starts with BEGIN word ("**+++idspispopd_____**" - 18 symbols)
- Then goes receiver (hash, name, address or some unique identifier = 32 symbols)
- Then goes type of payload ("T" = text, "B" = bytes)
- Then goes payload length (max 64*1024)
- Then goes message body (max 64*1024 symbols)
- Also, few service messages can be sent. They start with BEGIN word and have special hashes instead of receiver. Can have or not have body.


## Service tokens:

- **2hGXENUIwShosts1D3Z0W7l5KV4FqgYo** - get list of known lori hosts
- **oSeEtBW7MRsyncT4PD9Fg6idIbYXUfCn** - sync hosts (after token goes body with comma-separated hosts)
- **aoZJnr4pSEpingqhL6Vvb0dgXFM1sWwY** - ping receiver is in await list (after token goes receiver hash)

