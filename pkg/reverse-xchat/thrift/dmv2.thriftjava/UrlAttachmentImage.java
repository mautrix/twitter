package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000B\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0010\t\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\f\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\b\b\u0086\b\u0018\u0000 %2\u00020\u0001:\u0002&%B/\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0007¢\u0006\u0004\b\t\u0010\nJ\u0017\u0010\u000e\u001a\u00020\r2\u0006\u0010\f\u001a\u00020\u000bH\u0016¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J\u0012\u0010\u0012\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u0013J\u0012\u0010\u0014\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0014\u0010\u0011J\u0012\u0010\u0015\u001a\u0004\u0018\u00010\u0007HÆ\u0003¢\u0006\u0004\b\u0015\u0010\u0016J@\u0010\u0017\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u0007HÆ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u0010\u0010\u0019\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0019\u0010\u0011J\u0010\u0010\u001b\u001a\u00020\u001aHÖ\u0001¢\u0006\u0004\b\u001b\u0010\u001cJ\u001a\u0010 \u001a\u00020\u001f2\b\u0010\u001e\u001a\u0004\u0018\u00010\u001dHÖ\u0003¢\u0006\u0004\b \u0010!R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\"R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010#R\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\"R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00078\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010$¨\u0006'"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/UrlAttachmentImage;", "Lcom/bendb/thrifty/a;", "", "media_hash_key", "", "filesize_bytes", "filename", "Lcom/x/dmv2/thriftjava/MediaDimensions;", "dimensions", "<init>", "(Ljava/lang/String;Ljava/lang/Long;Ljava/lang/String;Lcom/x/dmv2/thriftjava/MediaDimensions;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Ljava/lang/Long;", "component3", "component4", "()Lcom/x/dmv2/thriftjava/MediaDimensions;", "copy", "(Ljava/lang/String;Ljava/lang/Long;Ljava/lang/String;Lcom/x/dmv2/thriftjava/MediaDimensions;)Lcom/x/dmv2/thriftjava/UrlAttachmentImage;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/lang/Long;", "Lcom/x/dmv2/thriftjava/MediaDimensions;", "Companion", "UrlAttachmentImageAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class UrlAttachmentImage implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final MediaDimensions dimensions;

    @JvmField
    @InterfaceC88465b
    public final String filename;

    @JvmField
    @InterfaceC88465b
    public final Long filesize_bytes;

    @JvmField
    @InterfaceC88465b
    public final String media_hash_key;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new UrlAttachmentImageAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/UrlAttachmentImage$UrlAttachmentImageAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/UrlAttachmentImage;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/UrlAttachmentImage;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/UrlAttachmentImage;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class UrlAttachmentImageAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public UrlAttachmentImage m85957read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            Long lValueOf = null;
            String string2 = null;
            MediaDimensions mediaDimensions = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new UrlAttachmentImage(string, lValueOf, string2, mediaDimensions);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            if (s != 4) {
                                C11272a.m14141a(protocol, b);
                            } else if (b == 12) {
                                mediaDimensions = (MediaDimensions) MediaDimensions.ADAPTER.read(protocol);
                            } else {
                                C11272a.m14141a(protocol, b);
                            }
                        } else if (b == 11) {
                            string2 = protocol.readString();
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 10) {
                        lValueOf = Long.valueOf(protocol.mo14124H0());
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 11) {
                    string = protocol.readString();
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a UrlAttachmentImage struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("UrlAttachmentImage");
            if (struct.media_hash_key != null) {
                protocol.mo14136v3("media_hash_key", 1, (byte) 11);
                protocol.mo14137w0(struct.media_hash_key);
            }
            if (struct.filesize_bytes != null) {
                protocol.mo14136v3("filesize_bytes", 2, (byte) 10);
                protocol.mo14121B3(struct.filesize_bytes.longValue());
            }
            if (struct.filename != null) {
                protocol.mo14136v3("filename", 3, (byte) 11);
                protocol.mo14137w0(struct.filename);
            }
            if (struct.dimensions != null) {
                protocol.mo14136v3("dimensions", 4, (byte) 12);
                MediaDimensions.ADAPTER.write(protocol, struct.dimensions);
            }
            protocol.mo14134i0();
        }
    }

    public UrlAttachmentImage(@InterfaceC88465b String str, @InterfaceC88465b Long l, @InterfaceC88465b String str2, @InterfaceC88465b MediaDimensions mediaDimensions) {
        this.media_hash_key = str;
        this.filesize_bytes = l;
        this.filename = str2;
        this.dimensions = mediaDimensions;
    }

    public static /* synthetic */ UrlAttachmentImage copy$default(UrlAttachmentImage urlAttachmentImage, String str, Long l, String str2, MediaDimensions mediaDimensions, int i, Object obj) {
        if ((i & 1) != 0) {
            str = urlAttachmentImage.media_hash_key;
        }
        if ((i & 2) != 0) {
            l = urlAttachmentImage.filesize_bytes;
        }
        if ((i & 4) != 0) {
            str2 = urlAttachmentImage.filename;
        }
        if ((i & 8) != 0) {
            mediaDimensions = urlAttachmentImage.dimensions;
        }
        return urlAttachmentImage.copy(str, l, str2, mediaDimensions);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getMedia_hash_key() {
        return this.media_hash_key;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Long getFilesize_bytes() {
        return this.filesize_bytes;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final String getFilename() {
        return this.filename;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final MediaDimensions getDimensions() {
        return this.dimensions;
    }

    @InterfaceC88464a
    public final UrlAttachmentImage copy(@InterfaceC88465b String media_hash_key, @InterfaceC88465b Long filesize_bytes, @InterfaceC88465b String filename, @InterfaceC88465b MediaDimensions dimensions) {
        return new UrlAttachmentImage(media_hash_key, filesize_bytes, filename, dimensions);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof UrlAttachmentImage)) {
            return false;
        }
        UrlAttachmentImage urlAttachmentImage = (UrlAttachmentImage) other;
        return Intrinsics.m65267c(this.media_hash_key, urlAttachmentImage.media_hash_key) && Intrinsics.m65267c(this.filesize_bytes, urlAttachmentImage.filesize_bytes) && Intrinsics.m65267c(this.filename, urlAttachmentImage.filename) && Intrinsics.m65267c(this.dimensions, urlAttachmentImage.dimensions);
    }

    public int hashCode() {
        String str = this.media_hash_key;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        Long l = this.filesize_bytes;
        int iHashCode2 = (iHashCode + (l == null ? 0 : l.hashCode())) * 31;
        String str2 = this.filename;
        int iHashCode3 = (iHashCode2 + (str2 == null ? 0 : str2.hashCode())) * 31;
        MediaDimensions mediaDimensions = this.dimensions;
        return iHashCode3 + (mediaDimensions != null ? mediaDimensions.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "UrlAttachmentImage(media_hash_key=" + this.media_hash_key + ", filesize_bytes=" + this.filesize_bytes + ", filename=" + this.filename + ", dimensions=" + this.dimensions + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
