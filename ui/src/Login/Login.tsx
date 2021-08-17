import { Button, Form, FormInstance, Input, Layout } from "antd";
import { Content } from "antd/lib/layout/layout";
import qs from 'qs'
import axios from 'axios';
import { openNotificationWithIcon } from "../notifications";
import { useHistory } from "react-router-dom";
import { useState } from "react";
import { useCookies } from "react-cookie";
import { AUTH_COOKIE_NAME } from "../Constants";
import { useAuthUserStorage } from "./UseUserStorage.hooks";

const validateMessages = {
  required: "${label} is required!",
};

export default function Login() {
  const [form] = Form.useForm()
  const history = useHistory();
  const [loading, setLoading] = useState(false);
  const [_, setAuthUser] = useAuthUserStorage();
  
  const onFinish = async (params) => {
    setLoading(true);
    try {
      const result = await axios({
        method: 'post',
        url: '/api/login',
        data: qs.stringify(params),
        headers: {
          'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
        }
      });
      setAuthUser(result.data);
      history.push('/');
    }
    catch (e) {
      setLoading(false);
      const validationMessage = (e?.response?.status == 401)
        ? 'The email or password was invalid'
        : 'An unexpected error occured.'
      form.setFields([
        {
          name: 'validationMessage',
          errors: [validationMessage],
        },
     ])
    }
  }
  
  return (
      <Layout>
        <Content className="centered-content">
          <div className="content-widget">
            <Form form={form}
              requiredMark={'optional'}
              name="nest-messages"
              onFinish={onFinish}
              validateMessages={validateMessages}
            >
              
              <Form.Item
                hasFeedback
                name={"email"}
                label="Email"
                rules={[{ required: true }]}
              >
                <Input />
              </Form.Item>
              <Form.Item
                hasFeedback
                name={"password"}
                rules={[{ required: true }]}
                label="Password"
              >
                <Input type="password" />
              </Form.Item>
              <Form.Item>
                <Button
                  loading={loading}
                  style={{ width: "100%" }}
                  type="primary"
                  htmlType="submit"
                >
                  Login
                </Button>
              </Form.Item>
              <Form.Item
               className="empty-form-item"
                name={"validationMessage"}
              >
              </Form.Item>
            </Form>
          </div>
        </Content>
      </Layout>
  );
}
