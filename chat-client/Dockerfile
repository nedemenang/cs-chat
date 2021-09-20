FROM node:14.17-alpine

WORKDIR /app

COPY package.json ./

# COPY package-lock.json ./

COPY ./ ./

# RUN npm cache verify

RUN npm install --no-package-lock

RUN chmod +x /app

CMD ["npm", "run", "start"]
