<template>
  <el-header class="navbar">
    <nav class="header">
      <router-link class="logo" to="/" tag="a">
        <img src="~/assets/images/logo.png" />
        <span>{{$t('subscan_essential')}}</span>
      </router-link>
      <div class="right-menu">
        <el-dropdown trigger="click">
          <div class="el-dropdown-link">
            {{$t('plugin')}}
            <i class="el-icon-arrow-down el-icon--right"></i>
          </div>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item>
              <router-link to="/blocks" tag="a">blocks</router-link>
            </el-dropdown-item>
            <el-dropdown-item>
              <router-link to="/events" tag="a">events</router-link>
            </el-dropdown-item>
            <el-dropdown-item
              v-for="item in pluginList"
              :key="item.name"
            >
              <router-link :to="`/${item.name}`" tag="a">{{item.name}}</router-link>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </nav>
  </el-header>
</template>
<script>
// import SearchInput from "@/views/Components/SearchInputNetwork";
// import { isSubscanHome } from "~/utils/tools";
import _ from "lodash";
export default {
  name: "NavBar",
  components: {
    // SearchInput,
  },
  props: {},
  data() {
    return {
      pluginList:[]
    };
  },
  async mounted() {
    let result = await this.$axios.$post("/api/scan/plugins");
    if (result.data && result.data.length > 0) {
      this.pluginList = _.filter(result.data,(item)=>{
        return item.ui;
      });
    }
  },
  computed: {},
  methods: {},
};
</script>
<style lang="scss" scoped>
@import "~assets/style/index.scss";
.navbar {
  display: flex;
  align-items: center;
  border-bottom: 1px solid #dcdfe6;
  .header {
    margin: 0 auto;
    width: 1180px;
    max-width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    .logo {
      img {
        height: 25px;
        width: 119px;
      }
      font-size: 12px;
      cursor: pointer;
    }
    .el-dropdown-link {
      cursor: pointer;
      line-height: 30px;
    }
  }
}
@media screen and (max-width: $screen-xs) {
}
</style>
<style>
li.el-dropdown-menu__item {
  list-style: none;
  line-height: 36px;
  padding: 0 20px;
  margin: 0;
  font-size: 14px;
  color: #606266;
  cursor: pointer;
  outline: none;
}
</style>
