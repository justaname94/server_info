import Search from "./pages/Search";
import History from "./pages/History";

export const routes = [
  {
    path: "",
    name: "home",
    components: {
      default: Search
    }
  },
  {
    path: "/history",
    name: "history",
    components: {
      default: History
    }
  }
];
