import * as amqp from 'amqplib';

class RabbitMQService {
  private conn: any;
  private channel: any;
  private exchange = 'events';

  async connect(url: string) {
    this.conn = await amqp.connect(url);
    this.channel = await this.conn.createChannel();
    await this.channel.assertExchange(this.exchange, 'topic', { durable: true });
  }

  async publish(routingKey: string, message: any, props: any = {}) {
    if (!this.channel) throw new Error('AMQP channel not initialized');
    const buf = Buffer.from(JSON.stringify(message));
    this.channel.publish(this.exchange, routingKey, buf, { persistent: true, ...props });
  }

  async subscribe(queueName: string, bindingKey: string, onMessage: (msg: any) => Promise<void>) {
    if (!this.channel) throw new Error('AMQP channel not initialized');
    await this.channel.assertQueue(queueName, { durable: true });
    await this.channel.bindQueue(queueName, this.exchange, bindingKey);
    await this.channel.consume(queueName, async (msg: any) => {
      if (!msg) return;
      try {
        const payload = JSON.parse(msg.content.toString());
        await onMessage(payload);
        this.channel.ack(msg);
      } catch (e) {
        console.error('Failed to process msg', e);
        this.channel.nack(msg, false, false);
      }
    });
  }
}

export const rabbit = new RabbitMQService();
