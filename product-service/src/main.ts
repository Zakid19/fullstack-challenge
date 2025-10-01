import 'reflect-metadata';
import {NestFactory} from '@nestjs/core'
import {AppModule} from './app.module';
import {RequestIdMiddleware} from './common/request-id.middleware';

async function bootstrap () {
  const app = await NestFactory.create(AppModule);
  app.use(new RequestIdMiddleware().use);
  await app.listen(3000);

  console.log('Product Service listening on 3000')
}

bootstrap();
