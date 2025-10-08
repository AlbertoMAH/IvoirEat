import { Row, Col, Card, Form, Input, Button, Typography } from "antd";

export const LoginPage = () => {
  return (
    <div style={{ padding: "50px" }}>
      <Row justify="center">
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card title="Login">
            <Typography.Title level={4}>Welcome Back!</Typography.Title>
            <Form
              layout="vertical"
              onFinish={(values) => console.log("Login values:", values)}
            >
              <Form.Item
                label="Email"
                name="email"
                rules={[{ required: true, message: "Please input your email!" }]}
              >
                <Input />
              </Form.Item>

              <Form.Item
                label="Password"
                name="password"
                rules={[{ required: true, message: "Please input your password!" }]}
              >
                <Input.Password />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit" block>
                  Log in
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  );
};
