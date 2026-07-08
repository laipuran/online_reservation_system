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
    route("services", "routes/services/page.tsx"),
    route("services/:id", "routes/services/service-detail.tsx"),
    route("my-reservations", "routes/_layout/my-reservations.tsx"),
    route("dashboard", "routes/dashboard.tsx"),
    route("complete-profile", "routes/_layout/complete-profile.tsx"),

    layout("routes/provider/_layout.tsx", [
      route("provider/services", "routes/provider/services/page.tsx"),
      route("provider/services/new", "routes/provider/services/new.tsx"),
      route("provider/services/:id", "routes/provider/services/service-detail.tsx"),
      route("provider/services/:id/edit", "routes/provider/services/edit.tsx"),
      route("provider/reservations", "routes/provider/reservations/page.tsx"),
    ]),
  ]),
] satisfies RouteConfig;
