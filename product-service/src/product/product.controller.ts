import { Controller, Post, Body, Get, Param, UsePipes, ValidationPipe } from '@nestjs/common';
import { ProductService } from './product.service';
import { CreateProductDto } from './dto/create-product.dto';

@Controller('products')
export class ProductController {
  constructor(private service: ProductService) {}

  @Post()
  @UsePipes(new ValidationPipe({ whitelist: true }))
  async create(@Body() dto: CreateProductDto) {
    return this.service.create(dto);
  }

  @Get(':id')
  async get(@Param('id') id: string) {
    return this.service.findById(id);
  }
}
