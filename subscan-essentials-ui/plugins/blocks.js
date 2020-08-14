export default {
  type: "page",
  body: {
    type: "crud",
    api: {
      method: "post",
      url: "api/scan/blocks",
      requestAdaptor: function (api) {
        return {
          ...api,
          data: {
            ...api.data,
            page: api.data.page - 1,
            row: api.data.perPage,
          },
        };
      },
      adaptor: function (payload, response) {
        return {
          ...payload,
          status: payload.code,
          data: {
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
        name: "block_num",
        label: "block_num",
      },
      {
        name: "finalized",
        label: "finalized",
      },
      {
        name: "extrinsics_count",
        label: "extrinsics_count",
      },
      {
        name: "event_count",
        label: "event_count",
      },
      {
        name: "validator",
        label: "validator",
      },
      {
        name: "hash",
        label: "hash",
      }
    ],
  },
}
