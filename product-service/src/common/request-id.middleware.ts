import {Request, Response, NextFunction} from 'express';
import {v4 as uuidv4 } from 'uuid'

export class RequestIdMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    const requestId = (req.headers['x-request-id'] as
      string) || uuidv4();
      (req as any).requestId = requestId;
      res.setHeader('x-request-id', requestId)
      next();
  }
}
