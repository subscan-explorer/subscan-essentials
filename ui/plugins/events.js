export default {
  type: "page",
  body: {
    type: "crud",
    api: {
      method: "post",
      url: "api/scan/events",
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
            items: payload.data.events,
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
        label: "block num",
      },
      {
        name: "event_idx",
        label: "event index",
      },
      {
        name: "extrinsic_hash",
        label: "extrinsic hash",
      },
      {
        name: "block_timestamp",
        label: "block timestamp",
      },
      {
        name: "module_id",
        label: "module",
      },
      {
        name: "event_id",
        label: "event",
      }
    ],
  },
}
