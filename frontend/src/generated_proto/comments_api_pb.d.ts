import * as $protobuf from "protobufjs";
import Long = require("long");
/** Namespace commentspb. */
export namespace commentspb {

    /** Properties of an AddComment. */
    interface IAddComment {

        /** AddComment id */
        id?: (string|null);

        /** AddComment senderId */
        senderId?: (string|null);

        /** AddComment itemId */
        itemId?: (string|null);

        /** AddComment itemType */
        itemType?: (string|null);

        /** AddComment content */
        content?: (string|null);

        /** AddComment categoryId */
        categoryId?: (string|null);

        /** AddComment parentId */
        parentId?: (string|null);
    }

    /** Represents an AddComment. */
    class AddComment implements IAddComment {

        /**
         * Constructs a new AddComment.
         * @param [properties] Properties to set
         */
        constructor(properties?: commentspb.IAddComment);

        /** AddComment id. */
        public id: string;

        /** AddComment senderId. */
        public senderId: string;

        /** AddComment itemId. */
        public itemId: string;

        /** AddComment itemType. */
        public itemType: string;

        /** AddComment content. */
        public content: string;

        /** AddComment categoryId. */
        public categoryId: string;

        /** AddComment parentId. */
        public parentId: string;

        /**
         * Creates a new AddComment instance using the specified properties.
         * @param [properties] Properties to set
         * @returns AddComment instance
         */
        public static create(properties?: commentspb.IAddComment): commentspb.AddComment;

        /**
         * Encodes the specified AddComment message. Does not implicitly {@link commentspb.AddComment.verify|verify} messages.
         * @param message AddComment message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: commentspb.IAddComment, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified AddComment message, length delimited. Does not implicitly {@link commentspb.AddComment.verify|verify} messages.
         * @param message AddComment message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: commentspb.IAddComment, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an AddComment message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns AddComment
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): commentspb.AddComment;

        /**
         * Decodes an AddComment message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns AddComment
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): commentspb.AddComment;

        /**
         * Verifies an AddComment message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an AddComment message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns AddComment
         */
        public static fromObject(object: { [k: string]: any }): commentspb.AddComment;

        /**
         * Creates a plain object from an AddComment message. Also converts values to other types if specified.
         * @param message AddComment
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: commentspb.AddComment, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this AddComment to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for AddComment
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }
}
