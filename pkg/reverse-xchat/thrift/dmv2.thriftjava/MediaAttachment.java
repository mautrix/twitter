package com.x.dmv2.thriftjava;

import android.gov.nist.core.C0003b;
import android.gov.nist.core.Separators;
import android.gov.nist.javax.sip.clientauthutils.C0026b;
import android.gov.nist.javax.sip.header.C0031b;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.ThriftException;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.twitter.app.p171di.app.C30722s1;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000F\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\t\n\u0002\b\b\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0012\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\t\b\u0086\b\u0018\u0000 22\u00020\u0001:\u000232Ba\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0006\u0012\b\u0010\t\u001a\u0004\u0018\u00010\b\u0012\b\u0010\n\u001a\u0004\u0018\u00010\b\u0012\b\u0010\u000b\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\f\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\r\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u000e\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u000f\u0010\u0010J\u0017\u0010\u0014\u001a\u00020\u00132\u0006\u0010\u0012\u001a\u00020\u0011H\u0016¢\u0006\u0004\b\u0014\u0010\u0015J\u0012\u0010\u0016\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0016\u0010\u0017J\u0012\u0010\u0018\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0018\u0010\u0019J\u0012\u0010\u001a\u001a\u0004\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b\u001a\u0010\u001bJ\u0012\u0010\u001c\u001a\u0004\u0018\u00010\bHÆ\u0003¢\u0006\u0004\b\u001c\u0010\u001dJ\u0012\u0010\u001e\u001a\u0004\u0018\u00010\bHÆ\u0003¢\u0006\u0004\b\u001e\u0010\u001dJ\u0012\u0010\u001f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001f\u0010\u0017J\u0012\u0010 \u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b \u0010\u0017J\u0012\u0010!\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b!\u0010\u0017J\u0012\u0010\"\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\"\u0010\u0017J|\u0010#\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u00062\n\b\u0002\u0010\t\u001a\u0004\u0018\u00010\b2\n\b\u0002\u0010\n\u001a\u0004\u0018\u00010\b2\n\b\u0002\u0010\u000b\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\f\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\r\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b#\u0010$J\u0010\u0010%\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b%\u0010\u0017J\u0010\u0010'\u001a\u00020&HÖ\u0001¢\u0006\u0004\b'\u0010(J\u001a\u0010,\u001a\u00020+2\b\u0010*\u001a\u0004\u0018\u00010)HÖ\u0003¢\u0006\u0004\b,\u0010-R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010.R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010/R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u00100R\u0016\u0010\t\u001a\u0004\u0018\u00010\b8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\t\u00101R\u0016\u0010\n\u001a\u0004\u0018\u00010\b8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\n\u00101R\u0016\u0010\u000b\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000b\u0010.R\u0016\u0010\f\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\f\u0010.R\u0016\u0010\r\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\r\u0010.R\u0016\u0010\u000e\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000e\u0010.¨\u00064"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MediaAttachment;", "Lcom/bendb/thrifty/a;", "", "media_hash_key", "Lcom/x/dmv2/thriftjava/MediaDimensions;", "dimensions", "Lcom/x/dmv2/thriftjava/MediaType;", "type", "", "duration_millis", "filesize_bytes", "filename", "attachment_id", "legacy_media_url_https", "legacy_media_preview_url", "<init>", "(Ljava/lang/String;Lcom/x/dmv2/thriftjava/MediaDimensions;Lcom/x/dmv2/thriftjava/MediaType;Ljava/lang/Long;Ljava/lang/Long;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Lcom/x/dmv2/thriftjava/MediaDimensions;", "component3", "()Lcom/x/dmv2/thriftjava/MediaType;", "component4", "()Ljava/lang/Long;", "component5", "component6", "component7", "component8", "component9", "copy", "(Ljava/lang/String;Lcom/x/dmv2/thriftjava/MediaDimensions;Lcom/x/dmv2/thriftjava/MediaType;Ljava/lang/Long;Ljava/lang/Long;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/MediaAttachment;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Lcom/x/dmv2/thriftjava/MediaDimensions;", "Lcom/x/dmv2/thriftjava/MediaType;", "Ljava/lang/Long;", "Companion", "MediaAttachmentAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MediaAttachment implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String attachment_id;

    @JvmField
    @InterfaceC88465b
    public final MediaDimensions dimensions;

    @JvmField
    @InterfaceC88465b
    public final Long duration_millis;

    @JvmField
    @InterfaceC88465b
    public final String filename;

    @JvmField
    @InterfaceC88465b
    public final Long filesize_bytes;

    @JvmField
    @InterfaceC88465b
    public final String legacy_media_preview_url;

    @JvmField
    @InterfaceC88465b
    public final String legacy_media_url_https;

    @JvmField
    @InterfaceC88465b
    public final String media_hash_key;

    @JvmField
    @InterfaceC88465b
    public final MediaType type;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MediaAttachmentAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MediaAttachment$MediaAttachmentAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MediaAttachment;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MediaAttachment;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MediaAttachment;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MediaAttachmentAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MediaAttachment m85947read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            MediaDimensions mediaDimensions = null;
            MediaType mediaType = null;
            Long lValueOf = null;
            Long lValueOf2 = null;
            String string2 = null;
            String string3 = null;
            String string4 = null;
            String string5 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b != 0) {
                    switch (c11265cMo14127V2.f38393b) {
                        case 1:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string = protocol.readString();
                                break;
                            }
                        case 2:
                            if (b != 12) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                mediaDimensions = (MediaDimensions) MediaDimensions.ADAPTER.read(protocol);
                                break;
                            }
                        case 3:
                            if (b != 8) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int iMo14132c4 = protocol.mo14132c4();
                                MediaType mediaTypeFindByValue = MediaType.INSTANCE.findByValue(iMo14132c4);
                                if (mediaTypeFindByValue == null) {
                                    throw new ThriftException(ThriftException.EnumC11260b.PROTOCOL_ERROR, C0031b.m45c(iMo14132c4, "Unexpected value for enum type MediaType: "));
                                }
                                mediaType = mediaTypeFindByValue;
                                break;
                            }
                        case 4:
                            if (b != 10) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                lValueOf = Long.valueOf(protocol.mo14124H0());
                                break;
                            }
                        case 5:
                            if (b != 10) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                lValueOf2 = Long.valueOf(protocol.mo14124H0());
                                break;
                            }
                        case 6:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string2 = protocol.readString();
                                break;
                            }
                        case 7:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string3 = protocol.readString();
                                break;
                            }
                        case 8:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string4 = protocol.readString();
                                break;
                            }
                        case 9:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string5 = protocol.readString();
                                break;
                            }
                        default:
                            C11272a.m14141a(protocol, b);
                            break;
                    }
                } else {
                    return new MediaAttachment(string, mediaDimensions, mediaType, lValueOf, lValueOf2, string2, string3, string4, string5);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MediaAttachment struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MediaAttachment");
            if (struct.media_hash_key != null) {
                protocol.mo14136v3("media_hash_key", 1, (byte) 11);
                protocol.mo14137w0(struct.media_hash_key);
            }
            if (struct.dimensions != null) {
                protocol.mo14136v3("dimensions", 2, (byte) 12);
                MediaDimensions.ADAPTER.write(protocol, struct.dimensions);
            }
            if (struct.type != null) {
                protocol.mo14136v3("type", 3, (byte) 8);
                protocol.mo14122C2(struct.type.value);
            }
            if (struct.duration_millis != null) {
                protocol.mo14136v3("duration_millis", 4, (byte) 10);
                protocol.mo14121B3(struct.duration_millis.longValue());
            }
            if (struct.filesize_bytes != null) {
                protocol.mo14136v3("filesize_bytes", 5, (byte) 10);
                protocol.mo14121B3(struct.filesize_bytes.longValue());
            }
            if (struct.filename != null) {
                protocol.mo14136v3("filename", 6, (byte) 11);
                protocol.mo14137w0(struct.filename);
            }
            if (struct.attachment_id != null) {
                protocol.mo14136v3("attachment_id", 7, (byte) 11);
                protocol.mo14137w0(struct.attachment_id);
            }
            if (struct.legacy_media_url_https != null) {
                protocol.mo14136v3("legacy_media_url_https", 8, (byte) 11);
                protocol.mo14137w0(struct.legacy_media_url_https);
            }
            if (struct.legacy_media_preview_url != null) {
                protocol.mo14136v3("legacy_media_preview_url", 9, (byte) 11);
                protocol.mo14137w0(struct.legacy_media_preview_url);
            }
            protocol.mo14134i0();
        }
    }

    public MediaAttachment(@InterfaceC88465b String str, @InterfaceC88465b MediaDimensions mediaDimensions, @InterfaceC88465b MediaType mediaType, @InterfaceC88465b Long l, @InterfaceC88465b Long l2, @InterfaceC88465b String str2, @InterfaceC88465b String str3, @InterfaceC88465b String str4, @InterfaceC88465b String str5) {
        this.media_hash_key = str;
        this.dimensions = mediaDimensions;
        this.type = mediaType;
        this.duration_millis = l;
        this.filesize_bytes = l2;
        this.filename = str2;
        this.attachment_id = str3;
        this.legacy_media_url_https = str4;
        this.legacy_media_preview_url = str5;
    }

    public static /* synthetic */ MediaAttachment copy$default(MediaAttachment mediaAttachment, String str, MediaDimensions mediaDimensions, MediaType mediaType, Long l, Long l2, String str2, String str3, String str4, String str5, int i, Object obj) {
        return mediaAttachment.copy((i & 1) != 0 ? mediaAttachment.media_hash_key : str, (i & 2) != 0 ? mediaAttachment.dimensions : mediaDimensions, (i & 4) != 0 ? mediaAttachment.type : mediaType, (i & 8) != 0 ? mediaAttachment.duration_millis : l, (i & 16) != 0 ? mediaAttachment.filesize_bytes : l2, (i & 32) != 0 ? mediaAttachment.filename : str2, (i & 64) != 0 ? mediaAttachment.attachment_id : str3, (i & 128) != 0 ? mediaAttachment.legacy_media_url_https : str4, (i & 256) != 0 ? mediaAttachment.legacy_media_preview_url : str5);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getMedia_hash_key() {
        return this.media_hash_key;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final MediaDimensions getDimensions() {
        return this.dimensions;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final MediaType getType() {
        return this.type;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final Long getDuration_millis() {
        return this.duration_millis;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final Long getFilesize_bytes() {
        return this.filesize_bytes;
    }

    @InterfaceC88465b
    /* renamed from: component6, reason: from getter */
    public final String getFilename() {
        return this.filename;
    }

    @InterfaceC88465b
    /* renamed from: component7, reason: from getter */
    public final String getAttachment_id() {
        return this.attachment_id;
    }

    @InterfaceC88465b
    /* renamed from: component8, reason: from getter */
    public final String getLegacy_media_url_https() {
        return this.legacy_media_url_https;
    }

    @InterfaceC88465b
    /* renamed from: component9, reason: from getter */
    public final String getLegacy_media_preview_url() {
        return this.legacy_media_preview_url;
    }

    @InterfaceC88464a
    public final MediaAttachment copy(@InterfaceC88465b String media_hash_key, @InterfaceC88465b MediaDimensions dimensions, @InterfaceC88465b MediaType type, @InterfaceC88465b Long duration_millis, @InterfaceC88465b Long filesize_bytes, @InterfaceC88465b String filename, @InterfaceC88465b String attachment_id, @InterfaceC88465b String legacy_media_url_https, @InterfaceC88465b String legacy_media_preview_url) {
        return new MediaAttachment(media_hash_key, dimensions, type, duration_millis, filesize_bytes, filename, attachment_id, legacy_media_url_https, legacy_media_preview_url);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MediaAttachment)) {
            return false;
        }
        MediaAttachment mediaAttachment = (MediaAttachment) other;
        return Intrinsics.m65267c(this.media_hash_key, mediaAttachment.media_hash_key) && Intrinsics.m65267c(this.dimensions, mediaAttachment.dimensions) && this.type == mediaAttachment.type && Intrinsics.m65267c(this.duration_millis, mediaAttachment.duration_millis) && Intrinsics.m65267c(this.filesize_bytes, mediaAttachment.filesize_bytes) && Intrinsics.m65267c(this.filename, mediaAttachment.filename) && Intrinsics.m65267c(this.attachment_id, mediaAttachment.attachment_id) && Intrinsics.m65267c(this.legacy_media_url_https, mediaAttachment.legacy_media_url_https) && Intrinsics.m65267c(this.legacy_media_preview_url, mediaAttachment.legacy_media_preview_url);
    }

    public int hashCode() {
        String str = this.media_hash_key;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        MediaDimensions mediaDimensions = this.dimensions;
        int iHashCode2 = (iHashCode + (mediaDimensions == null ? 0 : mediaDimensions.hashCode())) * 31;
        MediaType mediaType = this.type;
        int iHashCode3 = (iHashCode2 + (mediaType == null ? 0 : mediaType.hashCode())) * 31;
        Long l = this.duration_millis;
        int iHashCode4 = (iHashCode3 + (l == null ? 0 : l.hashCode())) * 31;
        Long l2 = this.filesize_bytes;
        int iHashCode5 = (iHashCode4 + (l2 == null ? 0 : l2.hashCode())) * 31;
        String str2 = this.filename;
        int iHashCode6 = (iHashCode5 + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.attachment_id;
        int iHashCode7 = (iHashCode6 + (str3 == null ? 0 : str3.hashCode())) * 31;
        String str4 = this.legacy_media_url_https;
        int iHashCode8 = (iHashCode7 + (str4 == null ? 0 : str4.hashCode())) * 31;
        String str5 = this.legacy_media_preview_url;
        return iHashCode8 + (str5 != null ? str5.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.media_hash_key;
        MediaDimensions mediaDimensions = this.dimensions;
        MediaType mediaType = this.type;
        Long l = this.duration_millis;
        Long l2 = this.filesize_bytes;
        String str2 = this.filename;
        String str3 = this.attachment_id;
        String str4 = this.legacy_media_url_https;
        String str5 = this.legacy_media_preview_url;
        StringBuilder sb = new StringBuilder("MediaAttachment(media_hash_key=");
        sb.append(str);
        sb.append(", dimensions=");
        sb.append(mediaDimensions);
        sb.append(", type=");
        sb.append(mediaType);
        sb.append(", duration_millis=");
        sb.append(l);
        sb.append(", filesize_bytes=");
        C30722s1.m39289b(l2, ", filename=", str2, ", attachment_id=", sb);
        C0026b.m37b(sb, str3, ", legacy_media_url_https=", str4, ", legacy_media_preview_url=");
        return C0003b.m4b(sb, str5, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
