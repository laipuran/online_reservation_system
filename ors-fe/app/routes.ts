import {
  type RouteConfig,
  index,
  route,
  layout,
} from "@react-router/dev/routes";

export default [
  layout("routes/_layout.tsx", [
    index("routes/home.tsx"),
    route("login", "routes/_layout/login.tsx"),
    route("register", "routes/_layout/register.tsx"),
    route("dashboard", "routes/dashboard.tsx"),
    route("complete-profile", "routes/_layout/complete-profile.tsx"),
    route("services/:id", "routes/services/service-detail.tsx"),
  ]),
] satisfies RouteConfig;
