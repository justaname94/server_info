import Vue from "vue";
import VueResource from "vue-resource";
import VueRouter from "vue-router";
import BootstrapVue from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import App from "./App.vue";
import { routes } from "./routes";

Vue.use(VueResource);
Vue.use(VueRouter);
Vue.use(BootstrapVue);

const router = new VueRouter({
  routes,
  mode: "history"
});

Vue.http.options.root = "http://localhost:8000/";

new Vue({
  el: "#app",
  router,
  render: h => h(App)
});
