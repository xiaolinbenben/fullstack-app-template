import { useMemo } from "react";
import { App as AntApp, ConfigProvider, Layout, theme } from "antd";
import { AppstoreOutlined } from "@ant-design/icons";
import { ProLayout } from "@ant-design/pro-components";
import OverviewPage from "./pages/OverviewPage";

export default function App() {
  const routes = useMemo(() => ({
    path: "/admin",
    routes: [{ path: "/admin/", name: "项目概览", icon: <AppstoreOutlined /> }],
  }), []);

  return (
    <ConfigProvider theme={{ algorithm: theme.defaultAlgorithm, token: { colorPrimary: "#1677ff", borderRadius: 6 } }}>
      <AntApp>
        <ProLayout
          title="Admin Console"
          logo={false}
          route={routes}
          location={{ pathname: window.location.pathname }}
          fixSiderbar
          menuItemRender={(item, dom) => <a href={item.path}>{dom}</a>}
          headerTitleRender={(logo, title) => <span className="header-title">{logo}{title}</span>}
        >
          <Layout.Content className="admin-content">
            <OverviewPage />
          </Layout.Content>
        </ProLayout>
      </AntApp>
    </ConfigProvider>
  );
}
