# 阶段一：构建前端产物
FROM node:18-alpine as build-stage
WORKDIR /app
# 如果遇到网络问题，可以加上 npm config set registry https://registry.npmmirror.com
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# 阶段二：Nginx 部署
FROM nginx:alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
# 替换默认配置为我们的自定义配置 (包含反向代理)
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]