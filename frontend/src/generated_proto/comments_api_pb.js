/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import * as $protobuf from "protobufjs/minimal";
// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;
// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});
export const commentspb = $root.commentspb = (() => {
    /**
     * Namespace commentspb.
     * @exports commentspb
     * @namespace
     */
    const commentspb = {};
    commentspb.AddComment = (function() {
        /**
         * Properties of an AddComment.
         * @memberof commentspb
         * @interface IAddComment
         * @property {string|null} [id] AddComment id
         * @property {string|null} [senderId] AddComment senderId
         * @property {string|null} [itemId] AddComment itemId
         * @property {string|null} [itemType] AddComment itemType
         * @property {string|null} [content] AddComment content
         * @property {string|null} [categoryId] AddComment categoryId
         * @property {string|null} [parentId] AddComment parentId
         */
        /**
         * Constructs a new AddComment.
         * @memberof commentspb
         * @classdesc Represents an AddComment.
         * @implements IAddComment
         * @constructor
         * @param {commentspb.IAddComment=} [properties] Properties to set
         */
        function AddComment(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * AddComment id.
         * @member {string} id
         * @memberof commentspb.AddComment
         * @instance
         */
        AddComment.prototype.id = "";
        /**
         * AddComment senderId.
         * @member {string} senderId
         * @memberof commentspb.AddComment
         * @instance
         */
        AddComment.prototype.senderId = "";
        /**
         * AddComment itemId.
         * @member {string} itemId
         * @memberof commentspb.AddComment
         * @instance
         */
        AddComment.prototype.itemId = "";
        /**
         * AddComment itemType.
         * @member {string} itemType
         * @memberof commentspb.AddComment
         * @instance
         */
        AddComment.prototype.itemType = "";
        /**
         * AddComment content.
         * @member {string} content
         * @memberof commentspb.AddComment
         * @instance
         */
        AddComment.prototype.content = "";
        /**
         * AddComment categoryId.
         * @member {string} categoryId
         * @memberof commentspb.AddComment
         * @instance
         */
        AddComment.prototype.categoryId = "";
        /**
         * AddComment parentId.
         * @member {string} parentId
         * @memberof commentspb.AddComment
         * @instance
         */
        AddComment.prototype.parentId = "";
        /**
         * Creates a new AddComment instance using the specified properties.
         * @function create
         * @memberof commentspb.AddComment
         * @static
         * @param {commentspb.IAddComment=} [properties] Properties to set
         * @returns {commentspb.AddComment} AddComment instance
         */
        AddComment.create = function create(properties) {
            return new AddComment(properties);
        };
        /**
         * Encodes the specified AddComment message. Does not implicitly {@link commentspb.AddComment.verify|verify} messages.
         * @function encode
         * @memberof commentspb.AddComment
         * @static
         * @param {commentspb.IAddComment} message AddComment message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AddComment.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.senderId != null && Object.hasOwnProperty.call(message, "senderId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.senderId);
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.itemId);
            if (message.itemType != null && Object.hasOwnProperty.call(message, "itemType"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.itemType);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.content);
            if (message.categoryId != null && Object.hasOwnProperty.call(message, "categoryId"))
                writer.uint32(/* id 6, wireType 2 =*/50).string(message.categoryId);
            if (message.parentId != null && Object.hasOwnProperty.call(message, "parentId"))
                writer.uint32(/* id 7, wireType 2 =*/58).string(message.parentId);
            return writer;
        };
        /**
         * Encodes the specified AddComment message, length delimited. Does not implicitly {@link commentspb.AddComment.verify|verify} messages.
         * @function encodeDelimited
         * @memberof commentspb.AddComment
         * @static
         * @param {commentspb.IAddComment} message AddComment message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AddComment.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes an AddComment message from the specified reader or buffer.
         * @function decode
         * @memberof commentspb.AddComment
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {commentspb.AddComment} AddComment
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AddComment.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.commentspb.AddComment();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                case 2: {
                        message.senderId = reader.string();
                        break;
                    }
                case 3: {
                        message.itemId = reader.string();
                        break;
                    }
                case 4: {
                        message.itemType = reader.string();
                        break;
                    }
                case 5: {
                        message.content = reader.string();
                        break;
                    }
                case 6: {
                        message.categoryId = reader.string();
                        break;
                    }
                case 7: {
                        message.parentId = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes an AddComment message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof commentspb.AddComment
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {commentspb.AddComment} AddComment
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AddComment.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies an AddComment message.
         * @function verify
         * @memberof commentspb.AddComment
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        AddComment.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                if (!$util.isString(message.senderId))
                    return "senderId: string expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.itemType != null && message.hasOwnProperty("itemType"))
                if (!$util.isString(message.itemType))
                    return "itemType: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            if (message.categoryId != null && message.hasOwnProperty("categoryId"))
                if (!$util.isString(message.categoryId))
                    return "categoryId: string expected";
            if (message.parentId != null && message.hasOwnProperty("parentId"))
                if (!$util.isString(message.parentId))
                    return "parentId: string expected";
            return null;
        };
        /**
         * Creates an AddComment message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof commentspb.AddComment
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {commentspb.AddComment} AddComment
         */
        AddComment.fromObject = function fromObject(object) {
            if (object instanceof $root.commentspb.AddComment)
                return object;
            let message = new $root.commentspb.AddComment();
            if (object.id != null)
                message.id = String(object.id);
            if (object.senderId != null)
                message.senderId = String(object.senderId);
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.itemType != null)
                message.itemType = String(object.itemType);
            if (object.content != null)
                message.content = String(object.content);
            if (object.categoryId != null)
                message.categoryId = String(object.categoryId);
            if (object.parentId != null)
                message.parentId = String(object.parentId);
            return message;
        };
        /**
         * Creates a plain object from an AddComment message. Also converts values to other types if specified.
         * @function toObject
         * @memberof commentspb.AddComment
         * @static
         * @param {commentspb.AddComment} message AddComment
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        AddComment.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.senderId = "";
                object.itemId = "";
                object.itemType = "";
                object.content = "";
                object.categoryId = "";
                object.parentId = "";
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                object.senderId = message.senderId;
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.itemType != null && message.hasOwnProperty("itemType"))
                object.itemType = message.itemType;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            if (message.categoryId != null && message.hasOwnProperty("categoryId"))
                object.categoryId = message.categoryId;
            if (message.parentId != null && message.hasOwnProperty("parentId"))
                object.parentId = message.parentId;
            return object;
        };
        /**
         * Converts this AddComment to JSON.
         * @function toJSON
         * @memberof commentspb.AddComment
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        AddComment.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for AddComment
         * @function getTypeUrl
         * @memberof commentspb.AddComment
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        AddComment.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/commentspb.AddComment";
        };
        return AddComment;
    })();
    return commentspb;
})();
export { $root as default };
