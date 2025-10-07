import { Row, Col, Card, Form, Input, Button, Typography } from "antd";

const { Title } = Typography;

export const LoginPage = () => {
  const onFinish = (values: any) => {
    console.log("Success:", values);
    // Ici, vous ajouteriez la logique d'authentification
  };

  return (
    <Row
      justify="center"
      align="middle"
      style={{
        height: "100vh",
        background: "#f0f2f5",
      }}
    >
      <Col xs={22} sm={16} md={12} lg={8} xl={6}>
        <Card>
          <div style={{ textAlign: "center", marginBottom: "24px" }}>
            <Title level={2}>Connexion</Title>
          </div>
          <Form
            name="login"
            onFinish={onFinish}
            layout="vertical"
            requiredMark="optional"
          >
            <Form.Item
              name="email"
              label="Email"
              rules={[{ required: true, type: "email", message: "Veuillez entrer un email valide!" }]}
            >
              <Input placeholder="Email" />
            </Form.Item>
            <Form.Item
              name="password"
              label="Mot de passe"
              rules={[{ required: true, message: "Veuillez entrer votre mot de passe!" }]}
            >
              <Input.Password placeholder="Mot de passe" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" block>
                Se connecter
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </Col>
    </Row>
  );
};
