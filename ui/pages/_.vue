<template>
  <div class="app-wrapper">
    <el-container direction="vertical">
      <nav-bar />
      <el-container direction="horizontal" class="main">
        <div id="amis"></div>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import NavBar from "~/components/NavBar.vue";
export default {
  name: "App",
  components: {
    NavBar,
  },
  data() {
    return {
    };
  },
  async mounted() {
    if (process.client) {
      let name = this.$route.params && this.$route.params.pathMatch;
      let result = await this.$axios.$post("/api/scan/plugins/ui", {name: name});
      if (result.data) {
        this.demo = result.data;
        this.initAmis();
      }
    }
  },
  methods: {
    initAmis() {
      var amis = window.amisRequire("amis/embed");
      amis.embed("#amis", this.demo);
    },
  },
};
</script>

<style lang='scss' scoped>
.app-wrapper {
  ::v-deep .a-Page {
    height: 100%;
  }
}
</style>
<style>
html,
body,
.app-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}
.app-wrapper {
  position: absolute;
  display: flex;
}
</style>
