/** Base class for all EventGraph errors. */
export class EventGraphError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "EventGraphError";
  }
}

export class OutOfRangeError extends EventGraphError {
  constructor(
    public readonly typeName: string,
    public readonly value: number,
    public readonly min: number,
    public readonly max: number,
  ) {
    super(`${typeName} value ${value} is out of range [${min}, ${max}]`);
    this.name = "OutOfRangeError";
  }
}

export class EmptyRequiredError extends EventGraphError {
  constructor(public readonly typeName: string) {
    super(`${typeName} cannot be empty`);
    this.name = "EmptyRequiredError";
  }
}

export class InvalidFormatError extends EventGraphError {
  constructor(
    public readonly typeName: string,
    public readonly providedValue: string,
    public readonly expectedFormat: string,
  ) {
    super(`${typeName} value "${providedValue}" does not match expected format: ${expectedFormat}`);
    this.name = "InvalidFormatError";
  }
}

export class InvalidTransitionError extends EventGraphError {
  constructor(
    public readonly from: string,
    public readonly to: string,
  ) {
    super(`Invalid transition from ${from} to ${to}`);
    this.name = "InvalidTransitionError";
  }
}

export class EventNotFoundError extends EventGraphError {
  constructor(public readonly eventId: string) {
    super(`Event ${eventId} not found`);
    this.name = "EventNotFoundError";
  }
}

export class ChainIntegrityError extends EventGraphError {
  constructor(
    public readonly position: number,
    detail: string,
  ) {
    super(`Chain integrity violation at position ${position}: ${detail}`);
    this.name = "ChainIntegrityError";
  }
}
