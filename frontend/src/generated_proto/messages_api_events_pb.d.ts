import * as $protobuf from "protobufjs";
import Long = require("long");
/** Namespace mes. */
export namespace mes {

    /** Properties of a SendMessage. */
    interface ISendMessage {

        /** SendMessage id */
        id?: (string|null);

        /** SendMessage conversationId */
        conversationId?: (string|null);

        /** SendMessage senderId */
        senderId?: (string|null);

        /** SendMessage recipientId */
        recipientId?: (string|null);

        /** SendMessage itemId */
        itemId?: (string|null);

        /** SendMessage body */
        body?: (string|null);

        /** SendMessage isRead */
        isRead?: (boolean|null);
    }

    /** Represents a SendMessage. */
    class SendMessage implements ISendMessage {

        /**
         * Constructs a new SendMessage.
         * @param [properties] Properties to set
         */
        constructor(properties?: mes.ISendMessage);

        /** SendMessage id. */
        public id: string;

        /** SendMessage conversationId. */
        public conversationId: string;

        /** SendMessage senderId. */
        public senderId: string;

        /** SendMessage recipientId. */
        public recipientId: string;

        /** SendMessage itemId. */
        public itemId: string;

        /** SendMessage body. */
        public body: string;

        /** SendMessage isRead. */
        public isRead: boolean;

        /**
         * Creates a new SendMessage instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SendMessage instance
         */
        public static create(properties?: mes.ISendMessage): mes.SendMessage;

        /**
         * Encodes the specified SendMessage message. Does not implicitly {@link mes.SendMessage.verify|verify} messages.
         * @param message SendMessage message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: mes.ISendMessage, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SendMessage message, length delimited. Does not implicitly {@link mes.SendMessage.verify|verify} messages.
         * @param message SendMessage message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: mes.ISendMessage, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SendMessage message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SendMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): mes.SendMessage;

        /**
         * Decodes a SendMessage message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SendMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): mes.SendMessage;

        /**
         * Verifies a SendMessage message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SendMessage message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SendMessage
         */
        public static fromObject(object: { [k: string]: any }): mes.SendMessage;

        /**
         * Creates a plain object from a SendMessage message. Also converts values to other types if specified.
         * @param message SendMessage
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: mes.SendMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SendMessage to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for SendMessage
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }
}
