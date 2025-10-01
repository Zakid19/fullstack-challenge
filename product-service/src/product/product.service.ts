import { Injectable, OnModuleInit } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Product } from './product.entity';
import { CreateProductDto } from './dto/create-product.dto';
import { redis } from '../common/redis.service';
import { rabbit } from '../common/rabbitmq.service';

@Injectable()
export class ProductService implements OnModuleInit {
  constructor(
    @InjectRepository(Product)
    private repo: Repository<Product>,
  ) {}

  async onModuleInit() {
    const url = process.env.RABBITMQ_URL || 'amqp://guest:guest@localhost:5672';
    try {
      await rabbit.connect(url);
      // subscribe to order.created events
      await rabbit.subscribe('product.order.listener', 'order.created', async (payload) => {
        console.log('product-service received order.created', payload);
        const productId = payload.productId;
        if (!productId) return;
        const product = await this.repo.findOneBy({ id: productId });
        if (!product) return;
        // assume 1 unit per order for demo
        product.qty = Math.max(0, product.qty - 1);
        await this.repo.save(product);
        await redis.del(`product:${productId}`);
      });
      console.log('Connected to RabbitMQ and subscribed to order.created');
    } catch (err) {
      console.error('Failed to connect to RabbitMQ', err);
    }
  }

  async create(dto: CreateProductDto) {
    const product = this.repo.create(dto as any);
    const saved = await this.repo.save(product);
    await rabbit.publish('product.created', saved);
    await redis.set(`product:${saved.id}`, JSON.stringify(saved), 'EX', 60);
    return saved;
  }

  async findById(id: string) {
    const key = `product:${id}`;
    const cached = await redis.get(key);
    if (cached) {
      try {
        return JSON.parse(cached);
      } catch {}
    }
    const product = await this.repo.findOneBy({ id });
    if (!product) return null;
    await redis.set(key, JSON.stringify(product), 'EX', 60);
    return product;
  }
}
