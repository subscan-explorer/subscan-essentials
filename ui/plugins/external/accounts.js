export default {
    type: "page",
    body: {
      type: "crud",
      api: {
        method: "post",
        url:"api/plugin/balance/accounts",
        requestAdaptor: function(api){
            return{
                ...api,
                data:{
                    ...api.data,
                    page: api.data.page -1,
                    row: api.data.perPage,
                },
            };
        },
        adapter: function(payload,response){
            return{
                ...payload,
                status: payload.code,
                data:{
                    items: payload.data.blocks,
                    count: payload.data.count
                },
                msg: payload.message
                };
            },
        },
        syncLocation: false,
        headerToolbar: [],
        columns: [{
            "name": "address",
            "label": "address"
          },
          {
            "name": "nonce",
            "label": "nonce"
          },
          {
            "name": "balance",
            "label": "balance"
          },
          {
            "name": "lock",
            "label": "lock"
          }
        ],
    },
}