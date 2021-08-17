import { notification } from "antd";

export const openNotificationWithIcon = (
  type: "success" | "error",
  title: string,
  description: string
) => {
  const notificationTypeFunc: Function = (notification as any)[type];
  notificationTypeFunc({
    message: title,
    description: description,
  });
};