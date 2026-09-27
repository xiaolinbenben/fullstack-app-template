import { useCallback, useEffect, useState } from "react";
import { CheckCircleOutlined, ReloadOutlined, WarningOutlined } from "@ant-design/icons";
import { Alert, Button, Card, Space, Tag, Typography } from "antd";
import { fetchHealth } from "../lib/api";

export default function OverviewPage() {
  const [state, setState] = useState<"loading" | "online" | "offline">("loading");
  const [message, setMessage] = useState("");

  const refresh = useCallback(async () => {
    setState("loading");
    try {
      const result = await fetchHealth();
      setState(result.data.status === "ok" ? "online" : "offline");
      setMessage("");
    } catch (error) {
      setState("offline");
      setMessage(error instanceof Error ? error.message : "服务不可用");
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return (
    <div className="overview-page">
      <div className="overview-heading">
        <div>
          <Typography.Text className="eyebrow">WORKSPACE / OVERVIEW</Typography.Text>
          <Typography.Title level={2}>项目概览</Typography.Title>
          <Typography.Paragraph type="secondary">管理端结构已经准备好，业务模块可以从这里开始扩展。</Typography.Paragraph>
        </div>
        <Button icon={<ReloadOutlined />} onClick={() => void refresh()} loading={state === "loading"}>刷新状态</Button>
      </div>

      <Card className="status-card" title="应用状态" extra={state === "online" ? <Tag color="success" icon={<CheckCircleOutlined />}>运行正常</Tag> : <Tag color="warning" icon={<WarningOutlined />}>{state === "loading" ? "检查中" : "暂不可用"}</Tag>}>
        <Space direction="vertical" size={4}>
          <Typography.Text>Gin 服务</Typography.Text>
          <Typography.Text type="secondary">GET /healthz</Typography.Text>
        </Space>
        {message && <Alert className="status-alert" type="warning" showIcon message={message} />}
      </Card>
    </div>
  );
}
