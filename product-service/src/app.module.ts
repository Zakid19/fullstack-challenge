import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Product } from './product/product.entity';
import { ProductModule } from './product/product.module';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: process.env.POSTGRES_HOST || 'localhost',
      port: +(process.env.POSTGRES_PORT || 5432),
      username: process.env.POSTGRES_USER || 'dev',
      password: process.env.POSTGRES_PASSWORD || 'dev',
      database: process.env.POSTGRES_DB || 'appdb',
      entities: [Product],
      synchronize: true,
    }),
    ProductModule,
  ],
})
export class AppModule {}
